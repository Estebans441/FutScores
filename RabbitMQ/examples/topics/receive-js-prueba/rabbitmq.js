import { connect } from "amqplib";

export default class RabbitMQClient {
    constructor(exchange, matches) {
        this.exchangeName = exchange;
        this.matches = matches;
    }

    async initiliaze() {
        this.connection = await connect("amqp://localhost:5672");
        this.channel = await this.connection.createChannel();
        this.exchange = await this.channel.assertExchange(
            this.exchangeName, 'topic',
            { 
                durable: true 
            }
        );
    }

    async subscribe() {
        await this.channel.assertQueue('', { exclusive: true }, (error2, q) => {
            if (error2) {
                throw error2;
            }
            this.queueName = q.queue;
        });


        this.matches.forEach((match) => {
            this.channel.bindQueue(
                this.queueName, 
                this.exchangeName, 
                match
            );
        });
        console.log("Waiting for messages...");

        this.channel.consume(this.queueName, (message) => {
            if (message != null) {
                console.log(`Received message: ${message.content.toString()}`);
            this.channel.ack(message);
            }
        });
    }
}
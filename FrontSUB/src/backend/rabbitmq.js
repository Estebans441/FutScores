import { connect } from "amqplib";

export default class RabbitMQClient {
    constructor(exchange, bindingKey, queue) {
        this.exchangeName = exchange;
        this.bindingKey = bindingKey;
        this.queueName = queue;
    }

    async initialize() {
        this.connection = await connect("amqp://localhost:5672");
        this.channel = await this.connection.createChannel();
        this.exchange = await this.channel.assertExchange(
            this.exchangeName, 
            "direct", 
            { 
                durable: true 
            }
        );
    }

    async subscribe(callback) {
        console.log("Waiting for messages...");
        await this.channel.assertQueue(this.queueName);
        await this.channel.bindQueue(
            this.queueName, 
            this.exchangeName, 
            this.bindingKey
        );

        this.channel.consume(this.queueName, (message) => {
            if (message != null) {
                console.log(`Received message: ${message.content.toString()}`);
                callback(message.content.toString());
                this.channel.ack(message);
            }
        });
    }
}
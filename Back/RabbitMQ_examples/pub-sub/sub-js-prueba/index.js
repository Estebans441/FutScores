import RabbitMQClient from './rabbitmq.js';

async function listen() {
    const rmqInstance = new RabbitMQClient("logs");
    await rmqInstance.initiliaze();
    await rmqInstance.subscribe();
}

listen();
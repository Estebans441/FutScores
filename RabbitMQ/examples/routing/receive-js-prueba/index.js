import RabbitMQClient from './rabbitmq.js';


var args = process.argv.slice(2);

if (args.length == 0) {
  console.log("Usage: receive_logs_direct.js [info] [warning] [error]");
  process.exit(1);
}

async function listen() {
    const rmqInstance = new RabbitMQClient("match_events", args);
    await rmqInstance.initiliaze();
    await rmqInstance.subscribe();
}

listen();
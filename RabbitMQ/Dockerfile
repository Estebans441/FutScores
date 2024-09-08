FROM rabbitmq:3.13-management

# Habilitar los plugins
RUN rabbitmq-plugins enable rabbitmq_stomp
RUN rabbitmq-plugins enable rabbitmq_web_stomp

EXPOSE 5672 15672 15674

CMD ["rabbitmq-server"]
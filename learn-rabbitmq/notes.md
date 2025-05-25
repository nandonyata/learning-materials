refers to this: https://adhon-rizky.medium.com/utilizing-rabbitmq-s-delay-capability-on-use-case-periodical-delayed-transaction-status-efb9c8ed3393



docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3-management

Then access the management UI at http://localhost:15672
(Default login: guest / guest)


//// WITH DELAYED PLUGINS ////

docker run -d --name rabbitmq \
  -p 5672:5672 -p 15672:15672 \
  -v $(pwd)/rabbitmq_delayed_message_exchange-3.13.0.ez:/plugins/rabbitmq_delayed_message_exchange-3.13.0.ez \
  -e RABBITMQ_PLUGINS_DIR="/plugins" \
  rabbitmq:3-management

docker exec -it docker_id /bin/bash
rabbitmq-plugins list
rabbitmq-plugins enable rabbitmq_delayed_message_exchange





//// ADJUSTED CREDS USER PASSWORD ////
docker run -d --name rabbitmq \
  -p 5672:5672 -p 15672:15672 \
  -v $(pwd)/rabbitmq_delayed_message_exchange-3.13.0.ez:/plugins/rabbitmq_delayed_message_exchange-3.13.0.ez \
  -e RABBITMQ_PLUGINS_DIR="/plugins" \
  -e RABBITMQ_DEFAULT_USER=myuser \
  -e RABBITMQ_DEFAULT_PASS=mysecurepassword \
  rabbitmq:3-management

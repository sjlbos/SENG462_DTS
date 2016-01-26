
using System;
using System.Collections.Generic;
using System.Configuration;
using System.Globalization;
using System.Linq;

namespace RabbitMessaging
{
    public class RabbitMessengerFactory
    {
        private static readonly IDictionary<string, HostElement> Hosts;
        private static readonly IDictionary<string, RabbitExcahangeElement> Exchanges;
        private static readonly IDictionary<string, RabbitQueueElement> Queues;
        private static readonly IDictionary<string, PublisherElement> Publishers; 
        private static readonly IDictionary<string, ReceiverElement> Receivers; 

        static RabbitMessengerFactory()
        {
            Hosts = new Dictionary<string, HostElement>();
            Exchanges = new Dictionary<string, RabbitExcahangeElement>();
            Queues = new Dictionary<string, RabbitQueueElement>();
            Publishers = new Dictionary<string, PublisherElement>();
            Receivers = new Dictionary<string, ReceiverElement>();

            var configuration = (RabbitConfigurationSection) ConfigurationManager.GetSection("rabbitmqConfig");

            foreach (HostElement host in configuration.Hosts)
            {
                Hosts.Add(host.Name, host);
            }

            foreach (RabbitExcahangeElement exchange in configuration.Exchanges)
            {
                Exchanges.Add(exchange.Name, exchange);
            }

            foreach (RabbitQueueElement queue in configuration.Queues)
            {
                Queues.Add(queue.Name, queue);
            }

            foreach (PublisherElement publisher in configuration.Publishers)
            {
                Publishers.Add(publisher.Name, publisher);
            }

            foreach (ReceiverElement receiver in configuration.Receivers)
            {
                Receivers.Add(receiver.Name, receiver);
            }
        }

        public static IMessagePublisher GetPublisher(string name)
        {
            if (name == null)
                throw new ArgumentNullException("name");
            if (!Publishers.ContainsKey(name))
                throw new RabbitConfigurationException(String.Format(CultureInfo.InvariantCulture,
                        "A publisher with name \"{0}\" could not be found in the configuration file.", name));

            var publisherConfig = Publishers[name];
            var hostConfig = GetHostConfiguration(publisherConfig.HostName);
            var exchangeConfig = GetExchangeConfiguration(publisherConfig.ExchangeName);
            QueueConfiguration deadLetterQueueConfig = null;
            if (publisherConfig.MessagesMustBeRouted)
            {
                if (String.IsNullOrWhiteSpace(publisherConfig.DeadLetterQueueName))
                {
                    throw new RabbitConfigurationException(String.Format(CultureInfo.InvariantCulture,
                        "Publisher \"{0}\" has \"MessagesMustBeRouted\" flag set to true but does not specify a dead letter queue.", name));
                }
                deadLetterQueueConfig = GetQueueConfiguration(publisherConfig.DeadLetterQueueName);
            }
            return new RabbitPublisher(hostConfig, exchangeConfig, publisherConfig.MessagesMustBeRouted, deadLetterQueueConfig);
        }

        public static IMessageReceiver GetReceiver(string name)
        {
            if (name == null)
                throw new ArgumentNullException("name");
            if (!Receivers.ContainsKey(name))
                throw new RabbitConfigurationException(String.Format(CultureInfo.InvariantCulture,
                    "A receiver with name \"{0}\" could not be found in the configuration file.", name));

            var receiverElement = Receivers[name];
            var hostConfig = GetHostConfiguration(receiverElement.HostName);
            var queueConfig = GetQueueConfiguration(receiverElement.QueueName);

            return new RabbitReceiver(hostConfig, queueConfig, receiverElement.AutoAckMessages, receiverElement.RequeueRejectedMessages);
        }

        private static HostConfiguration GetHostConfiguration(string hostName)
        {
            if(!Hosts.ContainsKey(hostName))
                throw new RabbitConfigurationException(String.Format(CultureInfo.InvariantCulture,
                    "A host with name \"{0}\" could not be found in the configuration file.", hostName));

            var hostElement = Hosts[hostName];
            return new HostConfiguration
            {
                HostName = hostElement.Name,
                VirtualHost = hostElement.VirtualHost,
                Port = hostElement.Port,
                Username = hostElement.Username,
                Password = hostElement.Password
            };
        }

        private static QueueConfiguration GetQueueConfiguration(string queueName)
        {
            if (!Queues.ContainsKey(queueName))
                throw new RabbitConfigurationException(String.Format(CultureInfo.InvariantCulture,
                    "A queue with name \"{0}\" could not be found in the configuration file.", queueName));
            
            var queueElement = Queues[queueName];
            var exchangeConfig = GetExchangeConfiguration(queueElement.ExchangeName);
            return new QueueConfiguration
            {
                Name = queueElement.Name,
                Exchange = exchangeConfig,
                AutoDelete = queueElement.AutoDelete,
                IsDurable = queueElement.IsDurable,
                IsExclusive = queueElement.IsExclusive,
                BindingKeys = GetBindingKeyListForQueue(queueElement)
            };
        }

        private static IList<string> GetBindingKeyListForQueue(RabbitQueueElement queueElement)
        {
            return (from BindingKeyElement key in queueElement.BindingKeys select key.Key).ToList();
        } 

        private static ExchangeConfiguration GetExchangeConfiguration(string exchangeName)
        {
            if(!Exchanges.ContainsKey(exchangeName))
                throw new RabbitConfigurationException(String.Format(CultureInfo.InvariantCulture,
                    "An exchange with name \"{0}\" could not be found in the configuration file.", exchangeName));

            var exchangeElement = Exchanges[exchangeName];
            return new ExchangeConfiguration
            {
                Name = exchangeElement.Name,
                ExchangeType = exchangeElement.Type,
                IsDurable = exchangeElement.IsDurable,
                AutoDelete = exchangeElement.AutoDelete
            };
        }
    }
}

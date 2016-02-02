
using System;
using System.Globalization;
using System.Text;
using RabbitMQ.Client;
using RabbitMQ.Client.Exceptions;

namespace RabbitMessaging
{
    public class RabbitPublisher : RabbitMessenger, IMessagePublisher
    {
        private readonly ExchangeConfiguration _exchangeConfig;
        private readonly bool _messagesMustBeRouted;
        private readonly QueueConfiguration _deadLetterQueueConfig;
        private IBasicProperties _sentMessageProperties;
            
        public RabbitPublisher(
            HostConfiguration hostConfig, 
            ExchangeConfiguration exchangeConfig, 
            bool messagesMustBeRouted,
            QueueConfiguration deadLetterQueueConfig
            ) 
        : base(hostConfig)
        {
            if (exchangeConfig == null)
                throw new ArgumentNullException("exchangeConfig");
            if (messagesMustBeRouted && deadLetterQueueConfig == null)
                throw new ArgumentNullException("deadLetterQueueConfig");

            _exchangeConfig = exchangeConfig;
            _messagesMustBeRouted = messagesMustBeRouted;
            _deadLetterQueueConfig = deadLetterQueueConfig;

            InitializeConnection();
        }

        private void InitializeConnection()
        {
            try
            {
                Channel = RabbitChannelProvider.OpenChannelToHost(HostConfig);
                CreateExchange(_exchangeConfig);
                _sentMessageProperties = Channel.CreateBasicProperties();
                _sentMessageProperties.ContentType = "text/plain";

                if (_messagesMustBeRouted)
                {
                    CreateQueue(_deadLetterQueueConfig);
                    _sentMessageProperties.DeliveryMode = 2; // 1 = Delivery is optional, 2 = Delivery is mandatory
                    Channel.BasicReturn += (sender, eventArgs) =>
                    {
                        Log.ErrorFormat(CultureInfo.InvariantCulture,
                            "Message sent by publisher could not be routed. Message: \"{0}\"",
                            Encoding.UTF8.GetString(eventArgs.Body));
                    };
                } 
            }
            catch (BrokerUnreachableException ex)
            {
                throw new ConnectionException("Publisher was unable to connect to broker. Check that the broker is online and the publisher is configured correctly.", ex);
            }
        }

        public void PublishMessage(string message, string routingKey)
        {
            if (message == null)
                throw new ArgumentNullException("message");

            if (Channel == null || Channel.IsClosed)
            {
                InitializeConnection();
            }

            byte[] messageBody = Encoding.UTF8.GetBytes(message);
            Channel.BasicPublish(_exchangeConfig.Name, routingKey, _sentMessageProperties, messageBody);
        }

        #region IDisposable

        public void Dispose()
        {
            Dispose(true);
            GC.SuppressFinalize(this);
        }

        protected virtual void Dispose(bool disposing)
        {
            if (disposing)
            {
                if (Channel != null)
                {
                    Channel.Close();
                    Channel = null;
                }
            }
        }

        #endregion
    }
}

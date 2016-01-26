
using System;
using System.Globalization;
using System.Reflection;
using log4net;
using RabbitMQ.Client;

namespace RabbitMessaging
{
    public abstract class RabbitMessenger
    {
        protected static readonly ILog Log = LogManager.GetLogger(MethodBase.GetCurrentMethod().DeclaringType);

        protected readonly HostConfiguration HostConfig;

        protected RabbitMessenger(HostConfiguration hostConfig)
        {
            if (hostConfig == null)
                throw new ArgumentNullException("hostConfig");

            HostConfig = hostConfig;
        }

        protected IModel Channel { get; set; }

        protected void CreateQueue(QueueConfiguration queueConfig)
        {
            if (queueConfig == null)
                throw new ArgumentNullException("queueConfig");
            if (Channel == null || Channel.IsClosed)
                throw new InvalidOperationException("Cannot create queue - channel is closed or null.");
            if (queueConfig.Exchange == null)
                throw new InvalidOperationException("Unable to create queue - queue configuration specifies null exchange configuration.");

            Log.DebugFormat(CultureInfo.InvariantCulture, "Creating exchange {0}.", queueConfig.Exchange.Name);
            CreateExchange(queueConfig.Exchange);

            Log.DebugFormat(CultureInfo.InvariantCulture, "Creating queue {0}.", queueConfig.Name);
            Channel.QueueDeclare(
                queueConfig.Name,
                queueConfig.IsDurable,
                queueConfig.IsExclusive,
                queueConfig.AutoDelete,
                null
                );

            BindQueueToExchange(queueConfig);
        }

        protected void CreateExchange(ExchangeConfiguration exchangeConfig)
        {
            if (exchangeConfig == null)
                throw new ArgumentNullException("exchangeConfig");
            if (Channel == null || Channel.IsClosed)
                throw new InvalidOperationException("Cannot create exchange - channel is closed or null.");

            Channel.ExchangeDeclare(
                exchangeConfig.Name, 
                exchangeConfig.ExchangeType,
                exchangeConfig.IsDurable,
                exchangeConfig.AutoDelete,
                null
                );
        }

        private void BindQueueToExchange(QueueConfiguration queueConfig)
        {
            foreach (var bindingKey in queueConfig.BindingKeys)
            {
                Log.DebugFormat(CultureInfo.InvariantCulture, "Binding queue {0} to exchange {1} with binding key {2}.", 
                    queueConfig.Name, 
                    queueConfig.Exchange.Name, 
                    bindingKey);
                Channel.QueueBind(queueConfig.Name, queueConfig.Exchange.Name, bindingKey);
            }
        }
    }
}

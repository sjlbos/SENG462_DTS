
using System;
using System.Configuration;

namespace RabbitMessaging
{
    public class RabbitConfigurationSection : ConfigurationSection
    {
        [ConfigurationProperty("hosts", IsRequired = true, IsDefaultCollection = true)]
        [ConfigurationCollection(typeof(HostElement), AddItemName = "host", CollectionType = ConfigurationElementCollectionType.BasicMap)]
        public HostCollection Hosts { get { return (HostCollection) this["hosts"]; } }

        [ConfigurationProperty("queues", IsRequired = true, IsDefaultCollection = true)]
        public RabbitQueueCollection Queues { get { return (RabbitQueueCollection) this["queues"]; } }

        [ConfigurationProperty("exchanges", IsRequired = true, IsDefaultCollection = true)]
        public RabbitExchangeCollection Exchanges { get { return (RabbitExchangeCollection) this["exchanges"]; } }

        [ConfigurationProperty("publishers", IsRequired = false, IsDefaultCollection = true)]
        public PublisherCollection Publishers { get { return (PublisherCollection) this["publishers"]; } }

        [ConfigurationProperty("receivers", IsRequired = false, IsDefaultCollection = true)]
        public ReceiverCollection Receivers { get { return (ReceiverCollection) this["receivers"]; } }
    }

    [ConfigurationCollection(typeof(PublisherCollection), AddItemName = "publisher", CollectionType = ConfigurationElementCollectionType.BasicMap)]
    public class PublisherCollection : ConfigurationElementCollection
    {
        protected override ConfigurationElement CreateNewElement()
        {
            return new PublisherElement();
        }

        protected override object GetElementKey(ConfigurationElement element)
        {
            return ((MessengerElement) element).Name;
        }
    }

    [ConfigurationCollection(typeof(ReceiverCollection), AddItemName = "receiver", CollectionType = ConfigurationElementCollectionType.BasicMap)]
    public class ReceiverCollection : ConfigurationElementCollection
    {
        protected override ConfigurationElement CreateNewElement()
        {
            return new ReceiverElement();
        }

        protected override object GetElementKey(ConfigurationElement element)
        {
            return ((MessengerElement) element).Name;
        }
    }

    public abstract class MessengerElement : ConfigurationElement
    {
        [ConfigurationProperty("name", IsRequired = true, IsKey = true)]
        public string Name
        {
            get { return (string) base["name"]; }
            set { base["name"] = value; }
        }

        [ConfigurationProperty("hostName", IsRequired = true)]
        public string HostName
        {
            get { return (string) base["hostName"]; }
            set { base["hostName"] = value; }
        }
    }

    public class PublisherElement : MessengerElement
    {
        [ConfigurationProperty("exchangeName", IsRequired = true)]
        public string ExchangeName
        {
            get { return (string) base["exchangeName"]; }
            set { base["exchangeName"] = value; }
        }

        [ConfigurationProperty("messagesMustBeRouted", IsRequired = false, DefaultValue = false)]
        public bool MessagesMustBeRouted
        {
            get { return (bool) base["messagesMustBeRouted"]; }
            set { base["messagesMustBeRouted"] = value; }
        }

        [ConfigurationProperty("deadLetterQueueName", IsRequired = false)]
        public string DeadLetterQueueName
        {
            get { return (string) base["deadLetterQueueName"]; }
            set { base["deadLetterQueueName"] = value; }
        }
    }

    public class ReceiverElement : MessengerElement
    {
        [ConfigurationProperty("queueName", IsRequired = true)]
        public string QueueName
        {
            get { return (string) base["queueName"]; }
            set { base["queueName"] = value; }
        }

        [ConfigurationProperty("autoAckMessages", IsRequired = false, DefaultValue = true)]
        public bool AutoAckMessages
        {
            get { return (bool) base["autoAckMessages"]; }
            set { base["autoAckMessages"] = value; }
        }

        [ConfigurationProperty("requeueRejectedMessages", IsRequired = false, DefaultValue = true)]
        public bool RequeueRejectedMessages
        {
            get { return (bool) base["requeueRejectedMessages"]; }
            set { base["requeueRejectedMessages"] = value; }
        }
    }

    public class HostCollection : ConfigurationElementCollection
    {
        protected override ConfigurationElement CreateNewElement()
        {
            return new HostElement();
        }

        protected override object GetElementKey(ConfigurationElement element)
        {
            return ((HostElement) element).Name;
        }
    }

    public class HostElement : ConfigurationElement
    {
        [ConfigurationProperty("name", IsRequired = true, IsKey = true)]
        public string Name
        {
            get { return (string) base["name"]; }
            set { base["name"] = value; }
        }

        [ConfigurationProperty("port", IsRequired = true)]
        public int Port
        {
            get { return (int) base["port"]; }
            set { base["port"] = value; }
        }

        [ConfigurationProperty("virtualHost", IsRequired = false)]
        public string VirtualHost
        {
            get { return (string) base["virtualHost"]; }
            set { base["virtualHost"] = value; }
        }

        [ConfigurationProperty("username", IsRequired = false)]
        public string Username
        {
            get { return (string) base["username"]; }
            set { base["username"] = value; }
        }

        [ConfigurationProperty("password", IsRequired = false)]
        public string Password
        {
            get { return (string) base["password"]; }
            set { base["password"] = value; }
        }
    }

    [ConfigurationCollection(typeof(RabbitQueueElement), AddItemName = "queue", CollectionType = ConfigurationElementCollectionType.BasicMap)]
    public class RabbitQueueCollection : ConfigurationElementCollection
    {
        protected override ConfigurationElement CreateNewElement()
        {
            return new RabbitQueueElement();
        }

        protected override object GetElementKey(ConfigurationElement element)
        {
            if (element == null)
                throw new ArgumentNullException("element");
            return ((RabbitQueueElement) element).Name;
        }
    }

    public class RabbitQueueElement : ConfigurationElement
    {
        [ConfigurationProperty("name", IsRequired = true, IsKey = true)]
        public string Name
        {
            get { return (string) base["name"]; }
            set { base["name"] = value; }
        }

        [ConfigurationProperty("exchangeName", IsRequired = true)]
        public string ExchangeName
        {
            get { return (string) base["exchangeName"]; }
            set { base["exchangeName"] = value; }
        }

        [ConfigurationProperty("isDurable", IsRequired = false, DefaultValue = true)]
        public bool IsDurable
        {
            get { return (bool) base["isDurable"]; }
            set { base["isDurable"] = value; }
        }

        [ConfigurationProperty("isExclusive", IsRequired = false, DefaultValue = false)]
        public bool IsExclusive
        {
            get { return (bool) base["isExclusive"]; }
            set { base["isExclusive"] = value; }
        }

        [ConfigurationProperty("autoDelete", IsRequired = false, DefaultValue = false)]
        public bool AutoDelete
        {
            get { return (bool) base["autoDelete"]; }
            set { base["autoDelete"] = value; }
        }

        [ConfigurationProperty("bindingKeys", IsRequired = false, IsDefaultCollection = true)]
        public BindingKeyCollection BindingKeys
        {
            get { return (BindingKeyCollection) base["bindingKeys"]; }
            set { base["bindingKeys"] = value; }
        }
    }

    public class BindingKeyCollection : ConfigurationElementCollection
    {
        protected override ConfigurationElement CreateNewElement()
        {
            return new BindingKeyElement();
        }

        protected override object GetElementKey(ConfigurationElement element)
        {   
            if(element == null)
                throw new ArgumentNullException("element");
            return ((BindingKeyElement) element).Key;
        }
    }

    public class BindingKeyElement : ConfigurationElement
    {
        [ConfigurationProperty("key", IsRequired = true, IsKey = true)]
        public string Key
        {
            get { return (string) base["key"]; }
            set { base["key"] = value; }
        }
    }

    [ConfigurationCollection(typeof(RabbitExcahangeElement), AddItemName = "exchange", CollectionType = ConfigurationElementCollectionType.BasicMap)]
    public class RabbitExchangeCollection : ConfigurationElementCollection
    {
        protected override ConfigurationElement CreateNewElement()
        {
            return new RabbitExcahangeElement();
        }

        protected override object GetElementKey(ConfigurationElement element)
        {
            if (element == null)
                throw new ArgumentNullException("element");
            return ((RabbitExcahangeElement) element).Name;
        }
    }

    public class RabbitExcahangeElement : ConfigurationElement
    {
        [ConfigurationProperty("name", IsRequired = true, IsKey = true)]
        public string Name
        {
            get { return (string) base["name"]; }
            set { base["name"] = value; }
        }

        [ConfigurationProperty("type", IsRequired = true)]
        public string Type
        {
            get { return (string) base["type"]; }
            set { base["type"] = value; }
        }

        [ConfigurationProperty("isDurable", IsRequired = false, DefaultValue = true)]
        public bool IsDurable
        {
            get { return (bool) base["isDurable"]; }
            set { base["isDurable"] = value; }
        }

        [ConfigurationProperty("autoDelete", IsRequired = false, DefaultValue = false)]
        public bool AutoDelete
        {
            get { return (bool) base["autoDelete"]; }
            set { base["autoDelete"] = value; }
        }
    }
}


using System;
using System.IO;
using System.Text;
using RabbitMQ.Client;
using RabbitMQ.Client.Exceptions;

namespace RabbitMessaging
{
    public class RabbitReceiver : RabbitMessenger, IMessageReceiver
    {
        private readonly QueueConfiguration _queueConfig;
        private readonly bool _autoAckMessages;
        private readonly bool _requeueRejectedMessages;

        private QueueingBasicConsumer _consumer;
        private string _consumerTag;

        public string LastMessage { get; private set; }
        private ulong _lastMessageTag;

        public RabbitReceiver(  HostConfiguration hostConfig, 
                                QueueConfiguration queueConfig, 
                                bool autoAckMessages = true , 
                                bool requeueRejectedMessages = true)
            :base(hostConfig)
        {
            if (queueConfig == null)
                throw new ArgumentNullException("queueConfig");

            _queueConfig = queueConfig;
            _autoAckMessages = autoAckMessages;
            _requeueRejectedMessages = requeueRejectedMessages;
            
            InitializeConnection();
        }

        private void InitializeConnection()
        {
            try
            {
                Channel = RabbitChannelProvider.OpenChannelToHost(HostConfig);
                CreateQueue(_queueConfig);
                _consumer = new QueueingBasicConsumer(Channel);
                _consumerTag = Channel.BasicConsume(_queueConfig.Name, !_autoAckMessages, _consumer);
            }
            catch (BrokerUnreachableException ex)
            {
                throw new ConnectionException("Receiver was unable to connect to broker. Check that the broker is online and the receiver is configured correctly.", ex);
            }
        }

        public string GetNextMessage()
        {
            if (Channel == null || Channel.IsClosed || _consumer == null)
            {
                InitializeConnection();
            }

            if (!_autoAckMessages && LastMessage != null)
            {
                throw new InvalidOperationException("The previous message must be acknowledged before a new message can be retrieved.");
            }

            try
            {
                var receiveEvent = _consumer.Queue.Dequeue();
                _lastMessageTag = receiveEvent.DeliveryTag;
                return Encoding.UTF8.GetString(receiveEvent.Body);
            }
            catch (EndOfStreamException)
            {
                Log.Info("Receiver was cancelled while waiting for messages.");
                return null;
            }
        }

        public void AckLastMessage()
        {
            if (_autoAckMessages)
                return;

            LastMessage = null;
            Channel.BasicAck(_lastMessageTag, false);
        }

        public void NackLastMessage()
        {
            if (_autoAckMessages)
                return;

            LastMessage = null;
            Channel.BasicNack(_lastMessageTag, false, _requeueRejectedMessages);
        }

        public void CancelMessageWait()
        {
            if (Channel != null)
            {
                Channel.BasicCancel(_consumerTag);
                Channel.Close();
            }
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
                CancelMessageWait();
                Channel = null;
            }
        }

        #endregion
    }
}

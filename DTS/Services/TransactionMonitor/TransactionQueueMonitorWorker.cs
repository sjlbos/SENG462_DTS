
using System;
using System.Collections.Concurrent;
using System.Globalization;
using log4net;
using Newtonsoft.Json;
using RabbitMessaging;
using ServiceHost;
using TransactionEvents;

namespace TransactionMonitor
{
    public class TransactionQueueMonitorWorker : QueueMonitorWorker
    {
        private static readonly ILog Log = LogManager.GetLogger(typeof (TransactionQueueMonitorWorker));

        private readonly BlockingCollection<TransactionEvent> _eventBuffer;
        private readonly TransactionEventConverter _deserializer;

        public TransactionQueueMonitorWorker(string instanceId, IMessageReceiver messageReceiver, BlockingCollection<TransactionEvent> eventBuffer) : base(instanceId, messageReceiver)
        {
            if(eventBuffer == null)
                throw new ArgumentNullException("eventBuffer");

            _eventBuffer = eventBuffer;
            _deserializer = new TransactionEventConverter();
        }

        public override void ProcessMessage(string message)
        {
            if (message == null)
                throw new ArgumentNullException("message");
            try
            {
                var transactionEvent = JsonConvert.DeserializeObject<TransactionEvent>(message, _deserializer);
                Log.DebugFormat(CultureInfo.InvariantCulture, "Worker {0} received transaction message with Id=\"{1}\".", InstanceId, transactionEvent.Id);
                _eventBuffer.Add(transactionEvent);
            }
            catch (UnrecognizedTransactionEventException ex)
            {
                Log.Error(String.Format(CultureInfo.InvariantCulture,
                    "Worker {0} received an unrecognized transaction event.", InstanceId), ex);
            }
        }
    }
}

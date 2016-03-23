
using System;
using System.Globalization;
using log4net;
using Newtonsoft.Json;
using RabbitMessaging;
using ServiceHost;
using TransactionEvents;
using TransactionMonitor.Repository;

namespace TransactionMonitor
{
    public class TransactionMonitorWorker : QueueMonitorWorker
    {
        private static readonly ILog Log = LogManager.GetLogger(typeof (TransactionMonitorWorker));

        private readonly TransactionEventConverter _deserializer;
        private readonly IAuditRepository _repository;

        public TransactionMonitorWorker(string instanceId, IMessageReceiver messageReceiver, IAuditRepository repository)
            : base(instanceId, messageReceiver)
        {
            if(repository == null)
                throw new ArgumentNullException("repository");

            _repository = repository;
            _deserializer = new TransactionEventConverter();
        }

        public override void ProcessMessage(string message)
        {
            if (message == null)
                throw new ArgumentNullException("message");
            try
            {
                var transactionEvent = JsonConvert.DeserializeObject<TransactionEvent>(message, _deserializer);
                Log.InfoFormat(CultureInfo.InvariantCulture, "Worker {0} received transaction message with Id=\"{1}\".", InstanceId, transactionEvent.Id);
                
                _repository.LogTransactionEvent(transactionEvent);
            }
            catch (UnrecognizedTransactionEventException ex)
            {
                Log.Error(String.Format(CultureInfo.InvariantCulture,
                    "Worker {0} received an unrecognized transaction event.", InstanceId), ex);
            }
        }
    }
}

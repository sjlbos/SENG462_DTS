
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
                
                WriteTransactionEventToRepository(transactionEvent);
            }
            catch (UnrecognizedTransactionEventException ex)
            {
                Log.Error(String.Format(CultureInfo.InvariantCulture,
                    "Worker {0} received an unrecognized transaction event.", InstanceId), ex);
            }
        }

        private void WriteTransactionEventToRepository(TransactionEvent transactionEvent)
        {
            try
            {
                if (transactionEvent is UserCommandEvent)
                {
                    _repository.LogUserCommandEvent(transactionEvent as UserCommandEvent);
                    return;
                }

                if (transactionEvent is QuoteServerEvent)
                {
                    _repository.LogQuoteServerEvent(transactionEvent as QuoteServerEvent);
                    return;
                }

                if (transactionEvent is AccountTransactionEvent)
                {
                    _repository.LogAccountTransactionEvent(transactionEvent as AccountTransactionEvent);
                    return;
                }

                if (transactionEvent is SystemEvent)
                {
                    _repository.LogSystemEvent(transactionEvent as SystemEvent);
                    return;
                }

                if (transactionEvent is ErrorEvent)
                {
                    _repository.LogErrorEvent(transactionEvent as ErrorEvent);
                    return;
                }

                if (transactionEvent is DebugEvent)
                {
                    _repository.LogDebugEvent(transactionEvent as DebugEvent);
                    return;
                }
            }
            catch (RepositoryException ex)
            {
                Log.Error(String.Format(CultureInfo.InvariantCulture,
                    "Worker {0} was unable to write transaction event with Id={1} to the repository.",
                    InstanceId, transactionEvent.Id), ex);
            }
        }
    }
}

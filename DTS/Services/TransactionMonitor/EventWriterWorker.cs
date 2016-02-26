using System;
using System.Collections.Concurrent;
using System.Globalization;
using System.Threading;
using log4net;
using ServiceHost;
using TransactionEvents;
using TransactionMonitor.Repository;

namespace TransactionMonitor
{
    public class EventWriterWorker : IWorker
    {
        private static readonly ILog Log = LogManager.GetLogger(typeof (EventWriterWorker));

        public string InstanceId { get; private set; }

        private readonly BlockingCollection<TransactionEvent> _eventBuffer;
        private readonly IAuditRepository _repository;

        public EventWriterWorker(string instanceId, IAuditRepository repository, BlockingCollection<TransactionEvent> eventBuffer)
        {
            if(instanceId == null)
                throw new ArgumentNullException("instanceId");
            if(repository == null)
                throw new ArgumentNullException("repository");
            if(eventBuffer == null)
                throw new ArgumentNullException("eventBuffer");

            InstanceId = instanceId;
            _eventBuffer = eventBuffer;
            _repository = repository;
        }

        public void Run(CancellationToken cancellationToken)
        {
            while (true)
            {
                if (cancellationToken.IsCancellationRequested)
                    break;

                try
                {
                    var receivedEvent = _eventBuffer.Take(cancellationToken);
                    WriteTransactionEventToRepository(receivedEvent);
                }
                catch (OperationCanceledException)
                {
                    break;
                }
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

        #region IDisposable

        public void Dispose()
        {
            Dispose(true);
            GC.SuppressFinalize(this);
        }

        protected virtual void Dispose(bool disposing)
        {

        }

        #endregion
    }
}


using System;
using System.Globalization;
using System.Threading;
using log4net;
using RabbitMessaging;

namespace ServiceHost
{
    public abstract class QueueMonitorWorker : IWorker
    {
        private static readonly ILog Log = LogManager.GetLogger(typeof(QueueMonitorWorker));

        public string InstanceId { get; private set; }
        protected readonly IMessageReceiver Receiver;

        protected QueueMonitorWorker(string instanceId, IMessageReceiver messageReceiver)
        {
            if (instanceId == null)
                throw new ArgumentNullException("instanceId");
            if (messageReceiver == null)
                throw new ArgumentNullException("messageReceiver");

            InstanceId = instanceId;
            Receiver = messageReceiver;
        }

        public void Run(CancellationToken cancellationToken)
        {
            if (cancellationToken.IsCancellationRequested)
                return;

            cancellationToken.Register(() => Receiver.CancelMessageWait());

            while (true)
            {
                if (cancellationToken.IsCancellationRequested)
                    break;

                string currentMessage = null;
                try
                {
                    currentMessage = Receiver.GetNextMessage();
                    if (cancellationToken.IsCancellationRequested)
                        break;
                    ProcessMessage(currentMessage, cancellationToken);
                    Receiver.AckLastMessage();
                }
                catch (Exception ex)
                {
                    Log.ErrorFormat(CultureInfo.InvariantCulture,
                        "Worker {0} encountered an error while attempting to process the message \"{1}\".",
                        InstanceId, currentMessage
                        );
                    Log.Error(ex);
                    Receiver.NackLastMessage();
                }
            }     
        }

        public abstract void ProcessMessage(string message, CancellationToken cancellationToken);

        #region IDisposable

        private bool _disposed = false;

        public void Dispose()
        {
            Dispose(true);
            GC.SuppressFinalize(this);
        }

        protected virtual void Dispose(bool disposing)
        {
            if (_disposed)
                return;

            if (disposing)
            {
                InstanceId = null;
            }

            if (Receiver != null)
            {
                Receiver.Dispose();
            }

            _disposed = true;
        }

        #endregion
    }
}

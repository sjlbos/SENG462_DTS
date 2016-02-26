
using System;
using System.Globalization;
using System.Reflection;
using System.Threading;
using log4net;

namespace ServiceHost
{
    public abstract class PollingWorker : IWorker
    {
        private static readonly ILog Log = LogManager.GetLogger(MethodBase.GetCurrentMethod().DeclaringType);

        public string InstanceId { get; private set; }
        public int PollRateMilliseconds { get; set; }

        protected PollingWorker(string id, int pollRateMilliseconds)
        {
            InstanceId = id;
            PollRateMilliseconds = pollRateMilliseconds;
        }

        public void Run(CancellationToken cancellationToken)
        {
            Log.Debug(String.Format(CultureInfo.InvariantCulture, "Worker {0} started.", InstanceId));
            // Perform an initial check to ensure worker was not shut down before starting.
            if (cancellationToken.IsCancellationRequested)
                return;

            while (true)
            {
                try
                {
                    DoWork();
                }
                catch (Exception ex)
                {
                    Log.Error(String.Format(CultureInfo.InvariantCulture, "Worker {0} encountered an unhandled exception.", InstanceId), ex);
                    throw;
                }

                // Check for a shutdown event after aquiring a lock
                if (cancellationToken.IsCancellationRequested)
                {
                    HandleShutdownEvent();
                    break;
                }

                Log.DebugFormat(CultureInfo.InvariantCulture, "Worker {0} sleeping...", InstanceId);
                var cancelled = cancellationToken.WaitHandle.WaitOne(PollRateMilliseconds);
                Log.DebugFormat(CultureInfo.InvariantCulture, "Worker {0} resuming...", InstanceId);
                    
                // Perform second shutdown check after resuming to avoid rentering the work loop.
                if (cancelled)
                {
                    HandleShutdownEvent();
                    break;
                }           
            }
            Log.DebugFormat(CultureInfo.InvariantCulture, "Worker {0} shutdown successfully.", InstanceId);
        }

        protected abstract void HandleShutdownEvent();
        protected abstract void DoWork();

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

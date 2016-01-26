using System;
using System.Collections.Generic;
using System.Globalization;
using System.Linq;
using System.Reflection;
using System.Threading;
using System.Threading.Tasks;
using log4net;

namespace ServiceHost
{
    public abstract class WorkerHost<T> : IService where T:IWorker
    {
        protected static readonly ILog Log = LogManager.GetLogger(MethodBase.GetCurrentMethod().DeclaringType);

        private IList<Task> _workerTasks; 
        private IList<T> _workers;   
        private readonly CancellationTokenSource _cancellationTokenSource;

        public bool RestartCrashedWorkers { get; set; }

        protected WorkerHost()
        {
            _cancellationTokenSource = new CancellationTokenSource();
            _workerTasks = new List<Task>();
        } 

        public void Start()
        {
            Log.Info("Starting service...");
            try
            {
                InitializeAndStartService();
            }
            catch (Exception ex)
            {
                Log.Fatal("The service encountered a fatal error and was unable to start.", ex);
                throw;
            }
            Log.Info("Service started successfully.");
        }

        public void Stop()
        {
            Log.Info("Stopping service...");
            _cancellationTokenSource.Cancel();
        }

        private void InitializeAndStartService()
        {
            InitializeService();
            _workers = GetWorkerList();
            if (_workers == null)
            {
                throw new ServiceException("Worker list returned by service was null.");
            }
            StartWorkers();
            Task.Run(() => MonitorWorkers());
        }

        private void StartWorkers()
        {
            _workerTasks = new Task[_workers.Count];
            for (int i = 0; i < _workers.Count; i ++)
            {
                var currentWorker = _workers[i]; // copy loop variable to local variable to avoid closure issues
                Log.Debug(String.Format(CultureInfo.InvariantCulture, "Starting worker {0}...", currentWorker.InstanceId));
                var workerTask = Task.Run(() => currentWorker.Run(_cancellationTokenSource.Token));
                _workerTasks[i] = workerTask;
            }
        }

        protected abstract void InitializeService();
        protected abstract IList<T> GetWorkerList();

        private void MonitorWorkers()
        {
            while (true)
            {
                if (_workerTasks.Count == 0)
                    break;

                int completedTaskIndex = Task.WaitAny(_workerTasks.ToArray());
                var completedTask = _workerTasks[completedTaskIndex];
                var completedWorker = _workers[completedTaskIndex];
                if (completedTask.IsFaulted)
                {
                    Log.Error(
                        String.Format(CultureInfo.InvariantCulture,
                            "Worker {0} encountered a fatal error and was stopped.", completedWorker.InstanceId),
                        completedTask.Exception);
                    if (RestartCrashedWorkers && !_cancellationTokenSource.IsCancellationRequested)
                    {
                        Log.InfoFormat(CultureInfo.InvariantCulture,
                            "Worker auto-restart enabled. Restarting worker {0}.", completedWorker.InstanceId);
                        _workerTasks[completedTaskIndex] =
                            Task.Run(() => completedWorker.Run(_cancellationTokenSource.Token));
                        continue;
                    }
                }
                else
                {
                    Log.InfoFormat(CultureInfo.InvariantCulture, "Worker {0} completed successfully.", completedWorker.InstanceId);
                } 
                _workerTasks.RemoveAt(completedTaskIndex);
                _workers.RemoveAt(completedTaskIndex);
            }
            Log.Info("All workers have shut down.");

            _cancellationTokenSource.Dispose();
        }
    }
}

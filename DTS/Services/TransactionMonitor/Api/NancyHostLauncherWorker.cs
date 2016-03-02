using System;
using System.Globalization;
using System.Threading;
using log4net;
using Nancy.Hosting.Self;
using ServiceHost;
using TransactionMonitor.Repository;

namespace TransactionMonitor.Api
{
    public class NancyHostLauncherWorker : IWorker
    {
        private static readonly ILog Log = LogManager.GetLogger(typeof (NancyHostLauncherWorker));

        public string InstanceId { get; private set; }

        private readonly Uri _baseUri;
        private readonly IAuditRepository _repository;

        public NancyHostLauncherWorker(string instanceId, Uri baseUri, IAuditRepository repository)
        {
            if(instanceId == null)
                throw new ArgumentNullException("instanceId");
            if(baseUri == null)
                throw new ArgumentNullException("baseUri");
            if (repository == null)
                throw new ArgumentNullException("repository");

            InstanceId = instanceId;
            _baseUri = baseUri;
            _repository = repository;
        }

        public void Run(CancellationToken cancellationToken)
        {
            if (cancellationToken.IsCancellationRequested)
                return;

            var nancyConfig = new HostConfiguration
            {
                UrlReservations = new UrlReservations
                {
                    CreateAutomatically = true
                }
            };

            var bootstrapper = new ApiBootstrapper(_repository);

            var host = new NancyHost(bootstrapper, nancyConfig, _baseUri);
            cancellationToken.Register(() =>
            {
                Log.InfoFormat("Stopping Nancy host at \"{0}\".", _baseUri);
                host.Stop();
                Log.InfoFormat("Nancy host at \"{0}\" successfully stopped.", _baseUri);
            });

            Log.InfoFormat(CultureInfo.InvariantCulture, "Starting Nancy host at \"{0}\".", _baseUri);
            host.Start();
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

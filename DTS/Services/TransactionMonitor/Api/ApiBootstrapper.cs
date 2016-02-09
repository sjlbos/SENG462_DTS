using System;
using System.Globalization;
using log4net;
using Nancy;
using Nancy.Bootstrapper;
using Nancy.TinyIoc;
using TransactionMonitor.Repository;

namespace TransactionMonitor.Api
{
    internal class ApiBootstrapper : DefaultNancyBootstrapper
    {
        private static readonly ILog Log = LogManager.GetLogger(typeof (AuditModule));

        private readonly IAuditRepository _repository;

        public ApiBootstrapper(IAuditRepository repository)
        {
            _repository = repository;
        }

        protected override void ConfigureApplicationContainer(TinyIoCContainer container)
        {
            base.ConfigureApplicationContainer(container);
            container.Register(_repository);
        }

        protected override void ApplicationStartup(TinyIoCContainer container, IPipelines pipelines)
        {
            base.ApplicationStartup(container, pipelines);

            pipelines.OnError += (context, ex) =>
            {
                Log.ErrorFormat(CultureInfo.InvariantCulture,
                    "An unhandled occured while executing the request \"{0} {1}{2}\".",
                    context.Request.Method,
                    context.Request.Path,
                    context.Request.Url.Query);
                Log.Error(ex);
                return null;
            };

            pipelines.BeforeRequest += (context) =>
            {
                Log.InfoFormat(CultureInfo.InvariantCulture,
                    "Received request \"{0} {1}{2}\".",
                    context.Request.Method,
                    context.Request.Path,
                    context.Request.Url.Query);
                return null;
            };

            pipelines.AfterRequest += (context) => Log.InfoFormat(CultureInfo.InvariantCulture,
                "Handled request \"{0} {1}{2}\". Response: ({3}) {4}",
                context.Request.Method,
                context.Request.Path,
                context.Request.Url.Query, 
                context.Response.StatusCode,
                context.Response.ReasonPhrase);
        }
    }
}

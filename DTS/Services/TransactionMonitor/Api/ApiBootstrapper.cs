using Nancy;
using TransactionMonitor.Repository;

namespace TransactionMonitor.Api
{
    internal class ApiBootstrapper : DefaultNancyBootstrapper
    {
        private readonly IAuditRepository _repository;

        public ApiBootstrapper(IAuditRepository repository)
        {
            _repository = repository;
        }

        protected override void ConfigureApplicationContainer(Nancy.TinyIoc.TinyIoCContainer container)
        {
            base.ConfigureApplicationContainer(container);
            container.Register(_repository);
        }
    }
}

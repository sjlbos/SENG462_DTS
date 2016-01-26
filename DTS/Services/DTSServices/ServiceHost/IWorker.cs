
using System.Threading;

namespace ServiceHost
{
    public interface IWorker
    {
        string InstanceId { get; }
        void Run(CancellationToken cancellationToken);
    }
}

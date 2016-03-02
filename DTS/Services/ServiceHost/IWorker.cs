
using System;
using System.Threading;

namespace ServiceHost
{
    public interface IWorker : IDisposable
    {
        string InstanceId { get; }
        void Run(CancellationToken cancellationToken);
    }
}

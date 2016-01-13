
using System;
using System.Timers;

namespace TopshelfService
{
    public class DemoService
    {
        readonly Timer _timer;

        public DemoService()
        {
            _timer = new Timer(1000) {AutoReset = true};
            _timer.Elapsed += (sender, eventArgs) => Console.WriteLine("It is {0} and all is well", DateTime.Now);
        }
        public void Start() { _timer.Start(); }
        public void Stop() { _timer.Stop(); }
    }
}

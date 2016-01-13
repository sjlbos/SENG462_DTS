using Topshelf;

namespace TopshelfService
{
    class Program
    {
        static void Main(string[] args)
        {
            HostFactory.Run(x =>                                 
            {
                x.UseLinuxIfAvailable();
                x.Service<Worker>(s =>    
                {
                    s.ConstructUsing(name => new Worker());    
                    s.WhenStarted(tc => tc.Start());             
                    s.WhenStopped(tc => tc.Stop());              
                });
                x.RunAsLocalSystem();                           

                x.SetDescription("Sample Topshelf Service");       
                x.SetDisplayName("TopshelfServicePrototype");                      
                x.SetServiceName("TopshelfServicePrototype");                       
            });                           
        }
    }
}

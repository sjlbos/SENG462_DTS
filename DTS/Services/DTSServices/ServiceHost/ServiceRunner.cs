
using System;
using Topshelf;

namespace ServiceHost
{
    public class ServiceRunner
    {
        /// <summary>
        /// Configures and launches a service of the specified type. The service will run under the local system account.
        /// </summary>
        /// <typeparam name="T">
        /// The type of service to run. Must be a class which implements the IService interface and contains a public default constructor.
        /// </typeparam>
        /// <param name="serviceName">The name of the service.</param>
        /// <param name="displayName">The name of the service as it will appear in the service manager.</param>
        /// <param name="description">The description of the service as it will appear in the service manager.</param>
        public static void RunService<T>(
            string serviceName, 
            string displayName, 
            string description) 
        where T:class, IService, new() 
        {
            if (String.IsNullOrWhiteSpace(serviceName)) 
                throw new ArgumentException("Parameter \"serviceName\" cannot be null or empty.");

            if (String.IsNullOrWhiteSpace(displayName))
            {
                displayName = serviceName;
            }

            if (description == null)
            {
                description = String.Empty;
            }

            HostFactory.Run(x =>
            {
                x.UseLinuxIfAvailable();
                x.Service<T>(s =>
                {
                    s.ConstructUsing(name => new T());
                    s.WhenStarted(tc => tc.Start());
                    s.WhenStopped(tc => tc.Stop());
                });
                x.RunAsLocalSystem();

                x.SetServiceName(serviceName);
                x.SetDisplayName(displayName);
                x.SetDescription(description);        
            });   
        }

        /// <summary>
        /// Configures and launches a service of the specified type. The service will run as the specified user.
        /// </summary>
        /// <typeparam name="T">
        /// The type of service to run. Must be a class which implements the IService interface and contains a public default constructor.
        /// </typeparam>
        /// <param name="serviceName">The name of the service.</param>
        /// <param name="displayName">The name of the service as it will appear in the service manager.</param>
        /// <param name="description">The description of the service as it will appear in the service manager.</param>
        /// <param name="username"></param>
        /// <param name="password"></param>
        public static void RunServiceAsUser<T>(
            string serviceName,
            string displayName,
            string description,
            string username,
            string password)
        where T : class, IService, new() 
        {
            if (String.IsNullOrWhiteSpace(serviceName))
                throw new ArgumentException("Parameter \"serviceName\" cannot be null or empty.");
            if (String.IsNullOrWhiteSpace(username))
                throw new ArgumentException("Parameter \"username\" cannot be null or empty.");
            if (password == null)
                throw new ArgumentNullException("password");

            if (String.IsNullOrWhiteSpace(displayName))
            {
                displayName = serviceName;
            }

            if (description == null)
            {
                description = String.Empty;
            }

            HostFactory.Run(x =>
            {
                x.UseLinuxIfAvailable();
                x.Service<T>(s =>
                {
                    s.ConstructUsing(name => new T());
                    s.WhenStarted(tc => tc.Start());
                    s.WhenStopped(tc => tc.Stop());
                });
                x.RunAs(username, password);

                x.SetServiceName(serviceName);
                x.SetDisplayName(displayName);
                x.SetDescription(description);
            });  
        }
    }
}

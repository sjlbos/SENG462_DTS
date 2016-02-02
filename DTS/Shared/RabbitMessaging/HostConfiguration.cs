
using System;

namespace RabbitMessaging
{
    public class HostConfiguration : IEquatable<HostConfiguration>
    {
        public static readonly int DefaultPort = 5672;
        public static readonly int DefaultTLSPort = 5671;

        public string HostName { get; set; }
        public string VirtualHost { get; set; }
        public int Port { get; set; }
        public string Username { get; set; }
        public string Password { get; set; }

        #region Equality Members

        public override int GetHashCode()
        {
            int hashCode = (HostName != null ? HostName.GetHashCode() : 0);
            hashCode = (hashCode*397) ^ (VirtualHost != null ? VirtualHost.GetHashCode() : 0);
            hashCode = (hashCode*397) ^ Port;
            hashCode = (hashCode*397) ^ (Username != null ? Username.GetHashCode() : 0);
            hashCode = (hashCode*397) ^ (Password != null ? Password.GetHashCode() : 0);
            return hashCode;
        }

        public override bool Equals(object obj)
        {
            return Equals(obj as HostConfiguration);
        }

        public bool Equals(HostConfiguration other)
        {
            if (other == null)
                return false;
            if (this == other)
                return true;
            return String.Equals(HostName, other.HostName) &&
                   String.Equals(VirtualHost, other.VirtualHost) &&
                   Port == other.Port &&
                   String.Equals(Username, other.Username) &&
                   String.Equals(Password, other.Password);
        }

        #endregion
    }
}

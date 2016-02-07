using System.Collections.Generic;
using System.IO;
using System.Xml;
using TransactionEvents;

namespace TransactionMonitor.Api
{
    internal static class XmlLog
    {
        public static void Write(IEnumerable<TransactionEvent> events, Stream stream)
        {
            using (var w = XmlWriter.Create(stream))
            {
                w.WriteStartDocument();
                w.WriteStartElement("log");
                foreach (var transactionEvent in events)
                {
                    transactionEvent.WriteXml(w);
                }
                w.WriteEndElement();
                w.WriteEndDocument(); 
            }
        }
    }
}

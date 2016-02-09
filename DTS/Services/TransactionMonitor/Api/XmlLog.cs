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
            var settings = new XmlWriterSettings
            {
                Indent = true,
                NewLineOnAttributes = true
            };
            using (var w = XmlWriter.Create(stream, settings))
            {
                w.WriteStartDocument();
                w.WriteStartElement("log");
                if (events != null)
                {
                    foreach (var transactionEvent in events)
                    {
                        transactionEvent.WriteXml(w);
                    }
                }
                w.WriteEndElement();
                w.WriteEndDocument(); 
            }
        }
    }
}

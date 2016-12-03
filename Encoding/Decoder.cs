using System.IO;
using System.Xml.Linq;

namespace Wsl.Encoding {
    public class Decoder {
        private XDocument xd;
        
        public Decoder() {
        }

        public Decoder(TextReader reader) {
            xd = XDocument.Load(reader);
        }

        public void ToWsl(TextWriter writer) {
            var transpiler = new XmlToWslTranspiler {
                DefaultNs = xd.Root.GetDefaultNamespace(),
                Indent = "",
                Writer = writer,
            };
            transpiler.WriteDocument(xd);
        }
    }
}
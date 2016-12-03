using System;
using System.IO;

namespace Wsl
{
    public class Program
    {
        public static int Main(string[] args)
        {
            var input = args.Length == 0 ? GetPipedInput() : ReadFile(args[0]);
            if (input == null) {
                Console.WriteLine("Missing document");
                return 1;
            }
            var decoder = new Encoding.Decoder(input);
            decoder.ToWsl(Console.Out);
            return 0;
        }

        private static TextReader GetPipedInput() {
            try {
                var isKey = Console.KeyAvailable;
                return null;
            } catch {
                return Console.In;
            }
        }

        private static StreamReader ReadFile(string filename) {
            var fi = new FileInfo(filename);
            if (!fi.Exists)
                return null;
            return fi.OpenText();
        }
    }
}

using System;
using System.Drawing;
using System.Drawing.Drawing2D;
using System.Drawing.Imaging;
using System.Runtime.InteropServices;

class MakeTransparent
{
    static void Main(string[] args)
    {
        if (args.Length < 2)
        {
            Console.WriteLine("Usage: MakeTransparent <input> <output> [outW] [outH] [colorThresh]");
            return;
        }

        string srcPath = args[0];
        string outPath = args[1];
        int outW = args.Length >= 3 ? int.Parse(args[2]) : 1920;
        int outH = args.Length >= 4 ? int.Parse(args[3]) : 1080;
        // Manhattan RGB distance to background; below this -> transparent
        int colorThresh = args.Length >= 5 ? int.Parse(args[4]) : 28;

        using (Bitmap src = new Bitmap(srcPath))
        {
            int w = src.Width;
            int h = src.Height;
            Console.WriteLine("Input: {0}x{1}", w, h);

            byte[] px = GetPacked(src);

            int bgB, bgG, bgR;
            DetectBackground(px, w, h, out bgB, out bgG, out bgR);
            Console.WriteLine("Background RGB: ({0},{1},{2})", bgR, bgG, bgB);

            // Soft fringe band: thresh .. softMax maps to alpha 0..255
            int softMax = Math.Max(colorThresh + 1, colorThresh * 4);

            for (int i = 0; i < w * h; i++)
            {
                int o = i * 4;
                int B = px[o], G = px[o + 1], R = px[o + 2];
                int dist = Math.Abs(B - bgB) + Math.Abs(G - bgG) + Math.Abs(R - bgR);

                byte a;
                if (dist <= colorThresh)
                {
                    a = 0;
                    // clear RGB to avoid fringe when composited on other backgrounds
                    px[o] = 0;
                    px[o + 1] = 0;
                    px[o + 2] = 0;
                }
                else if (dist < softMax)
                {
                    float t = (dist - colorThresh) / (float)(softMax - colorThresh);
                    if (t < 0f) t = 0f;
                    if (t > 1f) t = 1f;
                    a = (byte)(t * 255f + 0.5f);
                }
                else
                {
                    a = 255;
                }
                px[o + 3] = a;
            }

            // Find opaque content bounding box (with small padding)
            int minX = w, minY = h, maxX = -1, maxY = -1;
            for (int y = 0; y < h; y++)
            {
                for (int x = 0; x < w; x++)
                {
                    if (px[(y * w + x) * 4 + 3] < 8) continue;
                    if (x < minX) minX = x;
                    if (y < minY) minY = y;
                    if (x > maxX) maxX = x;
                    if (y > maxY) maxY = y;
                }
            }

            if (maxX < 0)
            {
                Console.WriteLine("No opaque content found.");
                return;
            }

            int pad = 4;
            minX = Math.Max(0, minX - pad);
            minY = Math.Max(0, minY - pad);
            maxX = Math.Min(w - 1, maxX + pad);
            maxY = Math.Min(h - 1, maxY + pad);
            int cw = maxX - minX + 1;
            int ch = maxY - minY + 1;
            Console.WriteLine("Content box: {0},{1} {2}x{3}", minX, minY, cw, ch);

            using (Bitmap content = new Bitmap(cw, ch, PixelFormat.Format32bppArgb))
            {
                byte[] crop = new byte[cw * ch * 4];
                for (int y = 0; y < ch; y++)
                {
                    for (int x = 0; x < cw; x++)
                    {
                        int si = ((minY + y) * w + (minX + x)) * 4;
                        int di = (y * cw + x) * 4;
                        crop[di] = px[si];
                        crop[di + 1] = px[si + 1];
                        crop[di + 2] = px[si + 2];
                        crop[di + 3] = px[si + 3];
                    }
                }
                SetPacked(content, crop);

                float margin = 0.08f;
                float maxW = outW * (1f - margin * 2f);
                float maxH = outH * (1f - margin * 2f);
                float scale = Math.Min(maxW / cw, maxH / ch);
                int dw = Math.Max(1, (int)Math.Round(cw * scale));
                int dh = Math.Max(1, (int)Math.Round(ch * scale));
                int dx = (outW - dw) / 2;
                int dy = (outH - dh) / 2;
                Console.WriteLine("Place: {0}x{1} at ({2},{3})", dw, dh, dx, dy);

                using (Bitmap final = new Bitmap(outW, outH, PixelFormat.Format32bppArgb))
                {
                    using (Graphics gClear = Graphics.FromImage(final))
                    {
                        gClear.Clear(Color.Transparent);
                    }

                    using (Graphics g = Graphics.FromImage(final))
                    {
                        g.CompositingMode = CompositingMode.SourceCopy;
                        g.CompositingQuality = CompositingQuality.HighQuality;
                        g.InterpolationMode = InterpolationMode.HighQualityBicubic;
                        g.SmoothingMode = SmoothingMode.HighQuality;
                        g.PixelOffsetMode = PixelOffsetMode.HighQuality;
                        g.DrawImage(content, new Rectangle(dx, dy, dw, dh),
                            new Rectangle(0, 0, cw, ch), GraphicsUnit.Pixel);
                    }

                    final.Save(outPath, ImageFormat.Png);
                    Console.WriteLine("Saved: {0} ({1}x{2})", outPath, final.Width, final.Height);
                }
            }
        }
    }

    // Sample four corners (+ nearby pixels) and take the median RGB as background.
    static void DetectBackground(byte[] px, int w, int h, out int bgB, out int bgG, out int bgR)
    {
        int[] xs = { 0, Math.Min(2, w - 1), Math.Max(0, w - 3), w - 1 };
        int[] ys = { 0, Math.Min(2, h - 1), Math.Max(0, h - 3), h - 1 };

        int n = xs.Length * ys.Length;
        int[] bs = new int[n];
        int[] gs = new int[n];
        int[] rs = new int[n];
        int k = 0;
        for (int yi = 0; yi < ys.Length; yi++)
        {
            for (int xi = 0; xi < xs.Length; xi++)
            {
                int o = (ys[yi] * w + xs[xi]) * 4;
                bs[k] = px[o];
                gs[k] = px[o + 1];
                rs[k] = px[o + 2];
                k++;
            }
        }

        Array.Sort(bs);
        Array.Sort(gs);
        Array.Sort(rs);
        int mid = n / 2;
        bgB = bs[mid];
        bgG = gs[mid];
        bgR = rs[mid];
    }

    static byte[] GetPacked(Bitmap bmp)
    {
        using (Bitmap tmp = new Bitmap(bmp.Width, bmp.Height, PixelFormat.Format32bppArgb))
        {
            using (Graphics g = Graphics.FromImage(tmp))
            {
                g.CompositingMode = CompositingMode.SourceCopy;
                g.DrawImage(bmp, 0, 0, bmp.Width, bmp.Height);
            }
            int w = tmp.Width, h = tmp.Height;
            BitmapData data = tmp.LockBits(new Rectangle(0, 0, w, h), ImageLockMode.ReadOnly, PixelFormat.Format32bppArgb);
            byte[] raw = new byte[data.Stride * h];
            Marshal.Copy(data.Scan0, raw, 0, raw.Length);
            int stride = data.Stride;
            tmp.UnlockBits(data);
            byte[] p = new byte[w * h * 4];
            for (int y = 0; y < h; y++)
                Buffer.BlockCopy(raw, y * stride, p, y * w * 4, w * 4);
            return p;
        }
    }

    static void SetPacked(Bitmap bmp, byte[] p)
    {
        int w = bmp.Width, h = bmp.Height;
        BitmapData data = bmp.LockBits(new Rectangle(0, 0, w, h), ImageLockMode.WriteOnly, PixelFormat.Format32bppArgb);
        byte[] raw = new byte[data.Stride * h];
        for (int y = 0; y < h; y++)
            Buffer.BlockCopy(p, y * w * 4, raw, y * data.Stride, w * 4);
        Marshal.Copy(raw, 0, data.Scan0, raw.Length);
        bmp.UnlockBits(data);
    }
}

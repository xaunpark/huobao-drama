import argparse
import math
import subprocess
import sys
from pathlib import Path

def create_contact_sheet(video_path: str, fps: float, out_dir: str, cols: int = 8):
    video = Path(video_path)
    if not video.exists():
        print(f"Error: Video file not found: {video_path}")
        sys.exit(1)

    print(f"Analyzing {video.name}...")
    # Lấy độ dài video bằng ffprobe
    probe = subprocess.run(
        ["ffprobe", "-v", "error", "-show_entries", "format=duration",
         "-of", "default=noprint_wrappers=1:nokey=1", str(video)],
        capture_output=True, text=True,
    )
    if probe.returncode != 0:
        print(f"Error running ffprobe: {probe.stderr}")
        sys.exit(1)
        
    try:
        duration = float(probe.stdout.strip())
    except ValueError:
        print(f"Error: Could not parse video duration from ffprobe: {probe.stdout}")
        sys.exit(1)
        
    total_frames = int(duration * fps)
    if total_frames == 0:
        print("Error: Video is too short or fps is too low.")
        sys.exit(1)
        
    rows = math.ceil(total_frames / cols)
    
    # Tên file đầu ra: TênVideo_contact_sheet.jpg
    out_path = Path(out_dir) / f"{video.stem}_contact_sheet.jpg"

    print(f"Generating contact sheet: {cols}x{rows} grid ({total_frames} frames at {fps}fps)")
    cmd = [
        "ffmpeg", "-y", "-i", str(video),
        "-vf", (
            f"fps={fps},"
            f"scale=320:-1,"
            f"drawtext=text='%{{pts\\:hms}}':x=5:y=5:fontsize=14:"
            f"fontcolor=white:borderw=1:bordercolor=black,"
            f"tile={cols}x{rows}"
        ),
        "-q:v", "2", str(out_path),
    ]
    
    result = subprocess.run(cmd, capture_output=True, text=True)
    if result.returncode != 0:
        print(f"Error generating contact sheet: {result.stderr}")
        sys.exit(1)
        
    print(f"Success! Saved to: {out_path}\n")

if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Generate a contact sheet for AI Video Review")
    parser.add_argument("video", help="Path to the video file")
    parser.add_argument("--fps", type=float, default=4.0, help="Frames per second to extract (default: 4.0 for Light mode, 8.0 for Deep mode)")
    parser.add_argument("--cols", type=int, default=8, help="Number of columns in the grid (default: 8)")
    parser.add_argument("--outdir", type=str, default="", help="Output directory (defaults to video's folder)")
    
    args = parser.parse_args()
    
    out_dir = args.outdir if args.outdir else str(Path(args.video).parent)
    create_contact_sheet(args.video, args.fps, out_dir, args.cols)

import ffmpeg
def convert(filename,output):
    stream = ffmpeg.input(filename)
    stream = ffmpeg.output(stream,output)
    ffmpeg.run_async(stream)

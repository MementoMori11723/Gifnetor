from PIL import Image, UnidentifiedImageError
import streamlit as st
import cv2
import os

def main():
    st.set_page_config(page_title="Gifnetor", page_icon=":vhs:", layout="centered")
    st.title("Video to Gif Converter")

    # Create the uploads directory if it doesn't exist
    if not os.path.exists("uploads"):
        os.makedirs("uploads")

    video = st.file_uploader("Upload Video", type=["mp4", "mov", "mkv"], accept_multiple_files=False)

    # Initialize `path` to None at the start
    path = None

    if video is not None:
        path = os.path.join("uploads/", video.name)

        try:
            with open(path, "wb") as upload:
                upload.write(video.getbuffer())
            st.success("Uploaded successfully")
        except Exception as e:
            st.error(f"Failed to upload video: {str(e)}")
            return

        button = st.button("Convert to Gif")

        if button:
            output, frame_count = convert(path)
            if output:
                st.image(output)
                st.info(f"Converted {frame_count} frames into a GIF")
                with open(output, "rb") as file:
                    st.download_button("Download Gif", data=file, file_name="output.gif", mime="image/gif")
            else:
                st.error("Failed to convert the video to GIF")

    reset = st.button("Remove Files")
    if reset and path:
        remove_files(path)

def convert(filepath):
    try:
        cap = cv2.VideoCapture(filepath)
        if not cap.isOpened():
            raise ValueError("Error occurred while opening the file")

        frames = []
        frame_count = 0

        while True:
            ret, frame = cap.read()
            if not ret:
                break
            try:
                pil_image = Image.fromarray(cv2.cvtColor(frame, cv2.COLOR_BGR2RGB))
                pil_image = pil_image.resize((pil_image.width // 4, pil_image.height // 4))
                frames.append(pil_image)
                frame_count += 1
            except UnidentifiedImageError as e:
                st.error(f"Error while processing frame: {str(e)}")
                cap.release()
                return None, 0

        if not frames:
            raise ValueError("No frames found in the video")

        gif_path = os.path.join("uploads/", "output.gif")
        frames[0].save(gif_path, save_all=True, append_images=frames[1:], loop=0, optimize=True, quality=50, duration=0)
        cap.release()

        return gif_path, frame_count

    except Exception as e:
        st.error(f"An error occurred during conversion: {str(e)}")
        return None, 0

def remove_files(path):
    try:
        output_gif = os.path.join("uploads/", "output.gif")
        if path and os.path.exists(path):
            os.remove(path)
        if os.path.exists(output_gif):
            os.remove(output_gif)
        st.warning("Files were removed") if path or output_gif else st.info("No files to remove")
    except Exception as e:
        st.error(f"An error occurred while removing files: {str(e)}")

if __name__ == "__main__":
    main()


from PIL import Image
import streamlit as st
import cv2
import os

def main():
    st.set_page_config(page_title="Gifnetor",page_icon=":vhs:",layout="centered")
    st.title("Video to Gif converter")
    video = st.file_uploader("Upload Video",type=["mp4","mov","mkv"],accept_multiple_files=False)
    path = os.path.join("uploads/",video.name) if video else None
    if video:
        with open(path,"wb") as upload :
            upload.write(video.getbuffer())
        st.success("Uploaded successfully")
        button = st.button("Convert to Gif")
        if button:
            output = convert(path)
            st.image(output)
            st.download_button("Download Gif",data=open(output,"rb"),file_name="output.gif",mime="image/gif")
    reset = st.button("Remove Files")
    if reset:
        if os.path.exists(path) and os.path.exists(os.path.join("uploads/","output.gif")):
            os.remove(path)
            os.remove(os.path.join("uploads/","output.gif"))
            st.warning(f"Files were removed")
        else:
            st.info(f"No Files were removed")
            
            
        

def convert(filepath):
    cap = cv2.VideoCapture(filepath)
    frames = []
    if not cap.isOpened():
        st.error("Error occured while opening the file")
        exit()
    while True:
        ret,frame = cap.read()
        if not ret:break
        pil_image = Image.fromarray(cv2.cvtColor(frame,cv2.COLOR_BGR2RGB))
        pil_image = pil_image.resize((pil_image.width // 4, pil_image.height // 4))
        frames.append(pil_image)
    gif_path = os.path.join("uploads/","output.gif")
    frames[0].save(gif_path,save_all=True,append_images=frames[1:],loop=0,optimize=True,quality=50,duration=0)
    cap.release()
    return gif_path

if __name__ == "__main__":
    main()
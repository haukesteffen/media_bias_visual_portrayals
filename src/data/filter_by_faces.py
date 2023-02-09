#!/usr/bin/env python

# ./filter_by_faces.py -src data/raw/ -out data/filtered/
# imports 
import argparse
import os
import shutil
import face_recognition

def get_flags():
    # setup argparser, add and parse arguments
    parser = argparse.ArgumentParser()
    parser.add_argument("-src", "--source", help="Source directory")
    parser.add_argument("-out", "--output", help="Output directory")
    args = parser.parse_args()
    src = args.source
    out = args.output

    # assert directories exist and are not the same
    assert os.path.isdir(src), print('Source directory does not exist')
    assert os.path.isdir(out), print('Output directory does not exist')
    assert src != out, print('Source and output can not be the same directory')

    return src, out

def list_images(dir):
    # returns a list of all images with .jpg type in input directory
    dir_list = os.listdir(dir)
    img_list = [img for img in dir_list if img.endswith('.jpg')]
    return img_list

def move_image(file, src, out):
    # moves file from src directory to out directory
    shutil.copy2(src+file, out+file)
    #os.remove(src+"/"+file)
    return

def get_face_encoding_from_image(img_path):
    # load image and create encoding
    image = face_recognition.load_image_file(img_path)
    encoding = face_recognition.face_encodings(image)
    return encoding

def check_for_faces(unknown_encoding, known_encodings):
    for enc in known_encodings:
        if face_recognition.compare_faces([enc], unknown_encoding[0])[0]:
            return True
    return False

def main():
    # get source and output directories from shell flags
    source_directory, output_directory = get_flags()

    # create list of files in source directory and filter for .jpg files
    img_list = list_images(source_directory)

    # get face encodings
    baerbock = get_face_encoding_from_image("data/training/annalena_baerbock.jpg")[0]
    laschet = get_face_encoding_from_image("data/training/armin_laschet.jpg")[0]
    scholz = get_face_encoding_from_image("data/training/olaf_scholz.jpg")[0]

    # move images to output directory if they contain a face
    for img in img_list:
        unknown_encoding = get_face_encoding_from_image(source_directory+img)
        if len(unknown_encoding) != 1:
            print(f"Multiple or no faces found in {img}")
            continue

        if check_for_faces(unknown_encoding, [baerbock, laschet, scholz]):
            move_image(img, source_directory, output_directory)
            print(f"Moved image {img}")
            continue

        print(f"No faces found for {img}")

    return 0

if __name__ == "__main__":
    main()
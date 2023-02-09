#!/usr/bin/env python

# ./filter_by_faces.py -src data/raw/ -out data/filtered/
# imports 
import argparse
import os
import shutil

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
    shutil.copy2(src+"/"+file, out+"/"+file)
    #os.remove(src+"/"+file)
    return

def check_for_faces(img):
    # todo
    return True

def main():
    # get source and output directories from shell flags
    source_directory, output_directory = get_flags()

    # create list of files in source directory and filter for .jpg files
    img_list = list_images(source_directory)

    # move images to output directory if they contain a face
    for img in img_list:
        if check_for_faces(img):
            move_image(img, source_directory, output_directory)

    print(f"Source directory {source_directory}\nOutput directory {output_directory}\nImages found in source directory: {img_list}")
    return 0

if __name__ == "__main__":
    main()
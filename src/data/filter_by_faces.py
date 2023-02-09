#!/usr/bin/env python

# ./filter_by_faces.py -src data/raw/ -out data/filtered/
# imports 
import argparse
import os

# setup argparser, add and parse arguments
parser = argparse.ArgumentParser()
parser.add_argument("-src", "--source", help="Source directory")
parser.add_argument("-out", "--output", help="Output directory")
args = parser.parse_args()
source_directory = args.source
output_directory = args.output

# assert directories exist and are not the same
assert os.path.isdir(source_directory), print('Source directory does not exist')
assert os.path.isdir(output_directory), print('Output directory does not exist')
assert source_directory != output_directory, print('Source and output can not be the same directory')

# create list of files in source directory and filter for .jpg files
dir_list = os.listdir(source_directory)
img_list = [img for img in dir_list if img.endswith('.jpg')]

print(f"Source directory {args.source}\nOutput directory {args.output}\nImages found in source directory: {img_list}")
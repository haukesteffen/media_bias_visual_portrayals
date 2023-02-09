# python src/data/filter_by_faces.py -src data/raw/ -out data/filtered/
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


dir_list = os.listdir(source_directory)


print(f"Source directory {args.source}\nOutput directory {args.output}\nFiles found in source directory: {dir_list}")
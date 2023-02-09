import argparse
parser = argparse.ArgumentParser()

# python src/data/filter_by_faces.py -src data/raw/ -out data/filtered/
parser.add_argument("-src", "--source", help="Source directory")
parser.add_argument("-out", "--output", help="Output directory")

args = parser.parse_args()

print(f"Source directory {args.source} Output directory {args.output}")
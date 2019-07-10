import json
import os


def filesPerProject():
    data = '[]'
    with open('filesPerProject.json', 'r') as myfile:
        data = myfile.read()

    d = json.loads(data)

    new = {}
    for x,y in d.iteritems():
        new[int(x)] = y

    with open("Output.json", "w") as text_file:
        text_file.write(json.dumps(new, sort_keys=True))



# ud = json.loads(js)
# print js # Unsorted JSON
# print ud # JSON to Dictionary (Unsorted)
# print json.dumps(ud, sort_keys=True) # This JSON will be sorted recursively

# from multiprocessing import Pool
# import multiprocessing
# import hashlib
# import os
# import re
# import string
# import boto3
# import subprocess


if __name__ == '__main__':
    filesPerProject()

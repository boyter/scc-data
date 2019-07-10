import json
import os


def filesPerProject():
    '''
    Converts the output of filesPerProject into something
    we can throw into a chart library since it needs to 
    be sorted
    '''
    data = '[]'
    with open('filesPerProject.json', 'r') as myfile:
        data = myfile.read()

    d = json.loads(data)

    new = []
    for x,y in d.iteritems():
        new.append([int(x), y])

    def cmp(a, b):
        if a == b:
            return 0
        if a < b:
            return -1
        return 1

    new.sort(cmp)

    with open("filesPerProject_converted.json", "w") as text_file:
        text_file.write(json.dumps(new, sort_keys=True))


if __name__ == '__main__':
    filesPerProject()

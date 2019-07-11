import json
import os


def filesPerProject():
    '''
    Converts the output of filesPerProject into something
    we can throw into a chart library since it needs to 
    be sorted
    It is a count of the number of projects that have a number of files

    EG. files:project where 123 projects have 2 files in them
    https://jsfiddle.net/uLw08scq/
    '''
    data = '[]'
    with open('./results/filesPerProject.json', 'r') as myfile:
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

    with open("./results/filesPerProject_converted.json", "w") as text_file:
        text_file.write(json.dumps(new, sort_keys=True))

def projectsPerLanguage():
    '''
    Converts output so we can see the number of projects per language
    https://jsfiddle.net/15v3c2pk/
    '''
    data = '[]'
    with open('./results/projectsPerLanguage.json', 'r') as myfile:
        data = myfile.read()

    d = json.loads(data)

    new = []
    for x,y in d.iteritems():
        new.append([x, y])

    def cmp(a, b):
        if a[1] == b[1]:
            return 0
        if a[1] > b[1]:
            return -1
        return 1

    new.sort(cmp)


    with open("./results/projectsPerLanguage_converted.json", "w") as text_file:
        text_file.write(json.dumps(new))


if __name__ == '__main__':
    filesPerProject()
    projectsPerLanguage()

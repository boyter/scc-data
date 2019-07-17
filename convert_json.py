'''
The job of this script is to convert the output of the Go job into 
what is required to fit whatever chart libary is being used and to 
possibly merge related sets together
'''

import json
import os
import operator


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
    for x, y in d.iteritems():
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


def mostCommonFileNames():
    '''
    Converts output so we can see the nmost common filenames
    '''
    data = '[]'
    with open('./results/fileNamesNoExtensionLowercaseCount.json', 'r') as myfile:
        data = myfile.read()

    d = json.loads(data)

    d = sorted(d.items(), key=operator.itemgetter(1), reverse=True)

    with open("./results/fileNamesNoExtensionLowercaseCount_converted.json", "w") as text_file:
        text_file.write(json.dumps(d))


def largestPerLanguage():
    '''
    Convert the largest files per language into markdown
    table for embedding
    '''
    data = '[]'
    with open('./results/largestPerLanguage.json', 'r') as myfile:
        data = myfile.read()

    d = json.loads(data)

    new = []
    for x, y in d.iteritems():
        new.append([x, y])

    def cmp(a, b):
        if a[1]['Bytes'] == b[1]['Bytes']:
            return 0
        if a[1]['Bytes'] > b[1]['Bytes']:
            return -1
        return 1

    new.sort(cmp)

    res = [
        '| language | filename | bytes |',
        '| -------- | -------- | ----- |',
    ]

    for y in new:
        x = '| %s | %s | %s |' % (y[0], y[1]['Filename'], y[1]['Bytes'])
        res.append(x)

    with open("./results/largestPerLanguage_converted.txt", "w") as text_file:
        text_file.write('''\n'''.join(res))


def longestPerLanguage():
    '''
    Convert the longest files per language into markdown
    table for embedding
    '''
    data = '[]'
    with open('./results/longestPerLanguage.json', 'r') as myfile:
        data = myfile.read()

    d = json.loads(data)

    new = []
    for x, y in d.iteritems():
        new.append([x, y])

    def cmp(a, b):
        if a[1]['Lines'] == b[1]['Lines']:
            return 0
        if a[1]['Lines'] > b[1]['Lines']:
            return -1
        return 1

    new.sort(cmp)

    res = [
        '| language | filename | lines |',
        '| -------- | -------- | ----- |',
    ]

    for y in new:
        x = '| %s | %s | %s |' % (y[0], y[1]['Filename'], y[1]['Lines'])
        res.append(x)

    with open("./results/longestPerLanguage_converted.txt", "w") as text_file:
        text_file.write('''\n'''.join(res))


def pureProjects():
    '''
    Converts the output of pureProjects into something
    we can throw into a chart library since it needs to 
    be sorted
    It is a count of the number of languages used by a project

    EG. languages:project where 123 projects have 2 languages in them
    https://jsfiddle.net/jqt81ufs/
    '''
    data = '[]'
    with open('./results/pureProjects.json', 'r') as myfile:
        data = myfile.read()

    d = json.loads(data)

    new = []
    for x, y in d.iteritems():
        new.append([int(x), y])

    def cmp(a, b):
        if a == b:
            return 0
        if a < b:
            return -1
        return 1

    new.sort(cmp)

    with open("./results/pureProjects_converted.json", "w") as text_file:
        text_file.write(json.dumps(new, sort_keys=True))


if __name__ == '__main__':
    filesPerProject()
    projectsPerLanguage()
    mostCommonFileNames()
    largestPerLanguage()
    longestPerLanguage()
    pureProjects()



'''
files
java 100
php 50
c 150

complexity
java 300
php 200
c 500


100 / 150 = 0.67
50 / 150 = 0.34
150 / 150 = 1


weighted complexity


java = 201
php = 68
c = 500
'''

'''
average complexity of Java repo
average complexity of Java repo between 1-50 files
average complexity of Java repo between 51-100 files
etc...
'''
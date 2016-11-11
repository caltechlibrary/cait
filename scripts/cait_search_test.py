#!/usr/bin/env python3
import urllib.request
import os
import sys
import json

class ServepagesTest:
    Version = 'v0.0.1'

    def __init__(self):
        self.jsonparse = json.JSONDecoder().decode
        self.site_url = os.environ.get('CAIT_SITE_URL')
        problem = False
        errors = []
        if not self.site_url:
            errors.append('Missing CAIT_SITE_URL in environment')
        if len(errors) > 0:
            raise Exception(', '.join(errors))

    def simple_search(self, f, q):
        '''simple_search(self, f, q) makes an http GET request to self.site_url passing f as the input name, and q as the value. Returns true if "</html>" found, false otherwise'''
        print('Check the logs for ' + self.site_url + ', and visually see if a "panic" is showing')
        data = urllib.parse.urlencode({f: q})
        data = data.encode('ascii')
        req = urllib.request.Request(self.site_url+'/search/basic/', data)
        with urllib.request.urlopen(req) as response:
            src = response.read().decode('UTF-8')
        print(src)
        if '</html>' not in src:
            return False
        return True


#
# Testing of CAIT's servepages
#
class Flags:
    flags = {}
    index = {}
    docs = []
    args = []
    parsed = {}

    def get(self, flag):
        key = flag.lstrip('-').strip()
        if flag in self.index:
            key = self.index[flag]
        if key in self.flags:
            return self.flags[key]
        return False


    def set(self, shortflag, longflag, default, msg):
        self.index[shortflag.strip()] = longflag.strip()
        self.flags[longflag.strip()] = default
        self.docs.append('-'+shortflag+', --'+longflag+'    '+msg)

    def printDefaults(self):
        for v in self.docs:
            print(v)

    def parse(self, args):
        for i in range(len(args)):
            arg = args[i]
            if '-' == arg[0:1]:
                key = arg.lstrip('-').strip()
                if key in self.index:
                    key = self.index[key]
                if key in self.flags:
                    val = self.flags[key]
                    if '=' in key:
                        (key, val) = key.split('=', 2)
                    elif (i+1) < len(args) and args[i+1][0:1] != '-':
                        val = args[i+1]
                    else:
                        val = True
                    self.flags[key] = val


class Testing:
    error_messages = []

    def error_count(self):
        return 'Total errors: {0}'.format(len(self.error_messages))

    def errors(self):
        return "\n".join(self.error_messages)

    def is_true(self, val, error_msg):
        if not val:
            self.error_messages.append(error_msg)
            return False
        return True

def main():
    '''Test ServepagesTest Class'''
    flags = Flags()
    flags.set('h', 'help', False, 'display help information')
    flags.set('v', 'version', False, 'display version information')
    flags.set('r', 'repository', False, 'display a list repositories')

    flags.parse(sys.argv)

    print('Checking for problems in environment...')
    try:
        servepages = ServepagesTest()
    except Exception as err:
        print('Unexpected error:', err)
        sys.exit(1)

    if flags.get('help'):
        print('USAGE %s [OPTIONS]' % sys.argv[0])
        flags.printDefaults()
        print('Version %s' % servepages.Version)
        sys.exit(0)

    if flags.get('version'):
        print('Version %s' % servepages.Version)
        sys.exit(0)

    t = Testing()
    print('Testing ' + servepages.site_url)
    t.is_true(servepages.simple_search("q", "Richard Feynman"), "testing ?q=Richard Feynman")
    try:
        t.is_true(servepages.simple_search("search", "Richarg Feynman"), "testing ?search=Richard Feynman")
    except Exception as err:
        print('Unexpected error: ', err)
    print(t.errors())
    print(t.error_count())


if __name__ == '__main__':
    main()

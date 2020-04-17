#!/usr/bin/python3

from datetime import datetime
import re
import subprocess
import sys

f = re.MULTILINE | re.IGNORECASE


def run_os_command(args):
    return subprocess.check_output(args, stderr=subprocess.DEVNULL)


def bazelifyScmVer(ver):
    # we expect '- example-v0.0.1-1-g1' as input, verify now
    wrong = re.sub('^- ([a-z0-9_-]+[a-z0-9_.-]+)$\n', '', ver, flags=f)
    if wrong != "":
        raise ValueError('Cannot parse this fragment: {}'.format(wrong))
    if ver == "":
        raise ValueError('Empty string')

    # lets make 'example-v0.0.1-1-g1'
    ver = re.sub(r'^- +', '', ver, flags=f)

    # lets make 'SCM_VERSION_EXAMPLE example-v0.0.1-1-g1'
    def scmUpper(m):
        return 'SCM_VERSION_'+m.groups()[0].upper()+" "+m.groups()[0]+m.groups()[1]
    ver = re.sub(r'^([a-z0-9_-]+)(-v[0-9].*)',
                 scmUpper, ver, flags=f)
    return ver


def bazelifySemVer(ver):
    # we expect 'SCM_VERSION_EXAMPLE example-v0.0.1-1-g1'
    # lets make 'SCM_VERSION_EXAMPLE example-v0.0.1-1+g1'
    ver = re.sub(r'^([ a-z0-9_-]+-v[0-9].*-[0-9]+)-(.*)',
                 r'\1+\2', ver, flags=f)

    # lets make 'SEM_VERSION_EXAMPLE 0.0.1-1+g1'
    ver = re.sub(r'^SCM(.* )([a-z0-9_-]+)-v([0-9].*)',
                 r'SEM\1\3', ver, flags=f)

    return ver

def bazelifySemRelVer(ver):

    # we expect 'SEM_VERSION_EXAMPLE 0.0.1-1+g1'
    # lets make 'SEM_RELEASE_VERSION_EXAMPLE 0.0.1-g1'
    ver = re.sub(r'^SEM([a-z0-9_-]+) ([0-9\.-].*)-[0-9]+\+(.*)',
                 r'SEM_RELEASE\1 \2-\3', ver, flags=f)

    return ver

# Since this script is so heavily used, lets make sure it gets unit tested.
# Temporarily, just unit test on *every* execution.
# TODO: Bazelify unit tests.

def UnitTestMe():
    tests = (
        {
            'desc': 'empty input',
            'input': '',
            'wantValueError': True,
        },
        {
            'desc': 'whitespace input',
            'input': ' \n\n\n\t\n',
            'wantValueError': True,
        },
        {
            'desc': 'bad input no hypen',
            'input': '- foo-v0.0.1\nbadexample\n',
            'wantValueError': True,
        },
        {
            'desc': 'usual input',
            'input': '- foo-v0.0.1\n- example-v99.99.99\n',
            'wantScm': 'SCM_VERSION_FOO foo-v0.0.1\nSCM_VERSION_EXAMPLE example-v99.99.99\n',
            'wantSem': 'SEM_VERSION_FOO 0.0.1\nSEM_VERSION_EXAMPLE 99.99.99\n',
        },
        {
            # TODO: This is the behavior of old bash script. Does it make sense?
            'desc': 'backward compatible - input with no version',
            'input': '- uninitialized\n',
            'wantScm': 'uninitialized\n',
            'wantSem': 'uninitialized\n',
        },
        {
            'desc': 'good input with git hash',
            'input': '- foo-v0.0.1-1-gdead4beef\n',
            'wantScm': 'SCM_VERSION_FOO foo-v0.0.1-1-gdead4beef\n',
            'wantSem': 'SEM_VERSION_FOO 0.0.1-1+gdead4beef\n',
        },
    )
    # We will loop through test cases now.
    # For each case we examine the result (what we got).
    # It must match expectations (what we want):
    #   - whether or not we got an exception for bad inputs
    #   - the SCM_VERSION output
    #   - the SEM_VERSION output
    for tt in tests:
        try:
            gotScm = bazelifyScmVer(tt['input'])
            if 'wantValueError' in tt and tt['wantValueError']:
                raise AssertionError('TestBazelifyPolyvers unit-test "{}" \ngot: no exception\nwant: ValueError exception'.format(
                    tt['desc']))

            if gotScm != tt['wantScm']:
                raise AssertionError('TestBazelifyPolyvers unit-test "{}" failed SCM\ngot: "{}"\nwant: "{}"'.format(
                    tt['desc'], gotScm, tt['wantScm']))

            gotSem = bazelifySemVer(gotScm)

            if gotSem != tt['wantSem']:
                raise AssertionError('TestBazelifyPolyvers unit-test "{}" failed SEM\ngot: "{}"\nwant: "{}"'.format(tt['desc'], gotSem, tt['wantSem']))

        except ValueError as e:
            if 'wantValueError' in tt and tt['wantValueError']:
                continue
            raise AssertionError('TestBazelifyPolyvers unit-test "{}"\ngot: exception {}\nwant: none'.format(tt['desc'], e))


######################## main #########################
if __name__ == "__main__":
    # UnitTestMe()

    # this will execute `polyvers status`
    git = run_os_command(['git', 'rev-parse', '--short', 'HEAD'])
    git = git.decode('utf-8')
    print('BUILD_SCM_REVISION {}'.format(git))

    git_branch = run_os_command(['git', 'rev-parse', '--abbrev-ref', 'HEAD'])
    git_branch = git_branch.decode('utf-8')
    print('BUILD_SCM_BRANCH {}'.format(git_branch))

    stamp = run_os_command(['git', 'show', '-s', '--format=%ct', 'HEAD'])
    stamp = int(stamp.decode('utf-8'))
    utctime = datetime.utcfromtimestamp(stamp)
    print(utctime.strftime('BUILD_SCM_TIMESTAMP %Y%m%dT%H%M%SZ'))

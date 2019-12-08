from behave import *
import subprocess
import time
import requests
import logging

formats = {
    "default": '$remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent" "$http_x_forwarded_for"',
    "upstream-time": '$remote_addr - $remote_user [$time_local] "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent" "$http_x_forwarded_for" $upstream_response_time'
}

bin_name = './prometheus-nginxlog-exporter'


@given('a running exporter listening on "{filename}" with format')
def run_exporter_impl(context, filename):
    filename = '.behave-sandbox/%s' % filename
    p = subprocess.Popen([bin_name, '--format', context.text, filename],
                         stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    time.sleep(.5)

    if p.returncode is not None:
        raise Exception('exporter exited too soon with exit code %d' % p.returncode)

    context.process = p


@given('a running exporter listening on "{filename}" with {format} format')
def run_exporter_impl(context, filename, format):
    filename = '.behave-sandbox/%s' % filename
    p = subprocess.Popen([bin_name, '--format', formats[format], filename],
                         stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    time.sleep(.5)

    if p.returncode is not None:
        raise Exception('exporter exited too soon with exit code %d' % p.returncode)

    context.process = p

@given(u'a running exporter listening with configuration file "{config}"')
def run_exporter_configfile_impl(context, config):
    p = subprocess.Popen([bin_name, '-config-file', 'features/%s' % config],
                         stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    time.sleep(.5)

    if p.returncode is not None:
        raise Exception('exporter exited too soon with exit code %d' % p.returncode)

    context.process = p

@when(u'the following HTTP request is logged to "{filename}"')
@when(u'the following HTTP requests are logged to "{filename}"')
def step_impl(context, filename):
    filename = '.behave-sandbox/%s' % filename
    with open(filename, 'a') as f:
        f.write("%s\n" % context.text)
    time.sleep(.5)


@when(u'the following HTTP request is logged to syslog on port {port}')
@when(u'the following HTTP requests are logged to syslog on port {port}')
def step_impl(context, port):
    log = logging.getLogger('test')
    log_handler = logging.handlers.SysLogHandler(("localhost", int(port)), logging.handlers.SysLogHandler.LOG_USER)
    log.addHandler(log_handler)

    lines = [l for l in context.text.split("\n") if l != ""]
    for l in lines:
        log.info(l)
    
    time.sleep(.5)


@then(u'the exporter should report value {val} for metric {metric}')
def step_impl(context, val, metric):
    while True:
        try:
            response = requests.get('http://localhost:4040/metrics')
            text = response.text
            break
        except requests.ConnectionError:
            continue

    lines = [l.strip("\n") for l in text.split("\n")]
    if not "%s %s" % (metric, val) in lines:
        raise AssertionError('expected metric "%s" could not be verified. Actual metrics:\n%s' % (context.text, text))

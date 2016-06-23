import os
import shutil


def before_scenario(context, scenario):
    if not os.path.isdir('.behave-sandbox'):
        os.mkdir('.behave-sandbox')

def after_scenario(context, scenario):
    shutil.rmtree('.behave-sandbox')
    if hasattr(context, 'process') and context.process.returncode is None:
        context.process.kill()

def after_step(context, step):
    if step.status == "failed" and hasattr(context, 'process'):
        # print("ohno failed")
        # print(context.process.returncode)
        if context.process.returncode is None:
            context.process.kill()
        out, err = context.process.communicate()

        print("stdout: %s\n\n" % out)
        print("stderr: %s\n\n" % err)

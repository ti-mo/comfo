from setuptools import setup

with open('../VERSION') as f:
    version = f.read().strip()

setup(name="comfo",
      version=version,
      description='Python API client for comfo, a Zehnder ComfoAir home automation controller.',
      license='MIT',
      py_modules=['comfo'],
      install_requires=[
        'twirp'
      ])

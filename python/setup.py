import setuptools

with open("../VERSION") as f:
    version = f.read().strip()

with open("../README.md", "r", encoding="utf-8") as fh:
    long_description = fh.read()

setuptools.setup(
    name="comfo",
    version=version,
    author="Timo Beckers",
    author_email="timo@incline.eu",
    description="Python API client for comfo, a Zehnder ComfoAir home \
        automation controller.",
    long_description=long_description,
    long_description_content_type="text/markdown",
    url="https://github.com/ti-mo/comfo",
    license="MIT",
    packages=setuptools.find_packages(),
    install_requires=["twirp==0.0.3"],
    python_requires=">=3.6",
)

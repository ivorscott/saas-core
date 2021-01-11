import setuptools

with open("README.md", "r") as fh:
    long_description = fh.read()

setuptools.setup(
    name="devpie-client-events",
    version="0.0.27",
    author="ivorscott",
    author_email="ivor@devpie.io",
    long_description=long_description,
    long_description_content_type="text/markdown",
    url="https://github.com/ivorscott/devpie-client-events",
    packages=setuptools.find_packages(),
    classifiers=[
        "Programming Language :: Python :: 3",
        "License :: OSI Approved :: MIT License",
        "Operating System :: OS Independent",
    ],
    python_requires='>=3.6',
)
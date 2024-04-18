# from distutils.core import setup
# from Cython.Build import cythonize

# setup(
#     ext_modules = cythonize("C:\\Users\\tyler\\Documents\\PROJECTS\\Dice_Analytics\\Dice_Analytics\\desktop\\dev\\server\\cython_modules\\run_cython.pyx")
# )

from setuptools import setup, Extension
from Cython.Build import cythonize
import numpy

# Fetch the numpy include directory.
numpy_include_dir = numpy.get_include()

extensions = [
    Extension("*", ["*.pyx"],
              include_dirs=[numpy_include_dir])  # Add include_dirs with numpy's include directory
]

setup(
    name='YourPackage',
    ext_modules=cythonize(extensions)
)
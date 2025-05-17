import math
from typing import Union, List


def convert_rotation_to_matrix(rotation: Union[float, int]) -> List[Union[float, int]]:
    rad = math.radians(rotation)
    return [math.cos(rad), -math.sin(rad), math.sin(rad), math.cos(rad)]

Program solving random kakurasu
===============================

A kakurasu is an n rows by m colums grid puzzle. The goal of the puzzle is to determine the black or 
white color of all cells by using the sums of weight of black cells in all rows and columns. The 
weight is a number from 1 to n. 
 
The image bellow shows a solved kakurasu. The numbers on the top and left sides are the row and column
weight. The numbers on the right and at the bottom are the sums of weight.

![Kakurasu example](Kakurasu_solution.jpg)


Solving algorithm
-----------------

Each sum has a limited set of possible weight combinations. From this set of possible solutions we can
deduce that some cells must be white and others must be black because they are respectively white or 
black in all solutions. The image below illustrate the deduction we can make from the sum 9. A cell
color is grey when it’s color is left unknown by the deduction.

![Deduction example](deduction.png)

Once we deduced the color of a cell, we can prune solutions with an incompatible color of the row or 
column contain the cell. By repeating the deduction and pruning operations, we can deduce the color 
of the grid cells. This deduction process ends when the color of all cells has been determined, or 
when no new deductions can be made. In the later case we are left with cells of unknowr color. This 
means that there are multiple solutions where the cells of unknown color are black and white.

To find the different solutions, we solve by assign the color white to a cell, and again by assigning 
the color black to that cell. This can be repeated as needed until the color of all cells has been 
determined for all solutions.
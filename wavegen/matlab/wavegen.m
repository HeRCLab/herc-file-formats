classdef wavegen
    properties (GetAccess=public,SetAccess=private)
        S = []
        T = []
        SampleRate
    end
    
    methods
        % constructor
        function self = wavegen(jsonfile)
            text = fileread(jsonfile);
            obj = jsondecode(text);
            self.S = obj.Signal.samples;
            self.T = obj.Signal.times;
            self.SampleRate = obj.Signal.SampleRate;
            if ~isequal(size(self.S), size(self.T))
                e = MException("MATLAB:error","Sample and time arrays are not the same length!");
                throw(e);
            end
   
        end
        
        % returns the number of samples in the file
        function Size(self)
            size(self.S, 1)
        end
    end
end